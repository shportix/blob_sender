package main

import (
	"bufio"
	"bytes"
	"fmt"
	depkeypair "gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/signcontrol"
	"io"
	"net/http"
	"os"
	"time"
)

const URL = "http://localhost:8010/blob/"

func CreateBlobRequest(data map[string]string, seedStr string) (*http.Request, error) {
	str := `{"data": {"id": 3,"data": {`
	num := 0
	for key := range data {
		num = 1
		str += `"` + key + `": "` + data[key] + `",`
	}
	if num != 0 {
		str = str[:len(str)-1]
	}
	str += `},"relationship": {"owner_id": ""}}}`
	var jsonStr = []byte(str)
	r, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return r, err
	}
	seed := depkeypair.MustParse(seedStr)
	err = signcontrol.SignRequest(r, seed)
	if err != nil {
		return r, err
	}
	return r, nil
}

func BlobListRequest(seedStr string) (*http.Request, error) {
	r, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return r, err
	}
	seed := depkeypair.MustParse(seedStr)
	err = signcontrol.SignRequest(r, seed)
	if err != nil {
		return r, err
	}
	return r, nil
}

func CheckTimeConstraint(seedStr string) (*http.Request, error) {
	r, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return r, err
	}
	r.Header.Set("date", time.Now().UTC().Add(-1*time.Hour).Format(http.TimeFormat))
	seed := depkeypair.MustParse(seedStr)
	err = signcontrol.SignRequest(r, seed)
	if err != nil {
		return r, err
	}
	return r, nil
}

func GetBlobByTdRequest(seedStr string, blobId string) (*http.Request, error) {
	r, err := http.NewRequest("GET", URL+blobId+"/", nil)
	if err != nil {
		return r, err
	}
	seed := depkeypair.MustParse(seedStr)
	err = signcontrol.SignRequest(r, seed)
	if err != nil {
		return r, err
	}
	return r, nil
}
func DeleteBlobByIdRequest(seedStr string, blobId string) (*http.Request, error) {
	url := "http://localhost:8010/blob/" + blobId + "/"
	fmt.Println(url)
	r, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return r, err
	}
	seed := depkeypair.MustParse(seedStr)
	err = signcontrol.SignRequest(r, seed)
	if err != nil {
		return r, err
	}
	return r, nil
}

func Do(r *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

func main() {
	e := false
	for !e {
		fmt.Println("Enter command letter:\ne-exit\nc-create\nd-delete\ng-get by id\nl - get list\nu - test unique\nt - time constraint test")
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		switch command[0] {
		case 'e':
			e = true
		case 'd':
			fmt.Println("Enter blob id:")
			blobId, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			blobId = blobId[:len(blobId)-2]
			fmt.Println("Enter your seed:")
			seedStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			seedStr = seedStr[:len(seedStr)-2]
			r, err := DeleteBlobByIdRequest(seedStr, blobId)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			Do(r)
		case 'l':
			fmt.Println("Enter your seed:")
			seedStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			seedStr = seedStr[:len(seedStr)-2]
			r, err := BlobListRequest(seedStr)
			if err != nil {
				fmt.Println(err)
				break
			}
			Do(r)
		case 'g':
			fmt.Println("Enter blob id:")
			blobId, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			blobId = blobId[:len(blobId)-2]
			fmt.Println("Enter your seed:")
			seedStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			seedStr = seedStr[:len(seedStr)-2]
			r, err := GetBlobByTdRequest(seedStr, blobId)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			Do(r)
		case 'c':
			s := false
			data := make(map[string]string)
			for !s {
				fmt.Println("Enter key, or s to stop adding data:")
				key, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("An error occured while reading input. Please try again", err)
					return
				}
				switch key[:len(key)-2] {
				case "s":
					s = true
				default:
					fmt.Println("Enter value:")
					value, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("An error occured while reading input. Please try again", err)
						return
					}
					data[key[:len(key)-2]] = value[:len(value)-2]
				}
			}
			fmt.Println("Enter your seed:")
			seedStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			seedStr = seedStr[:len(seedStr)-2]
			r, err := CreateBlobRequest(data, seedStr)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			Do(r)
		case 'u':
			fmt.Println("Enter your seed:")
			seedStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			seedStr = seedStr[:len(seedStr)-2]
			r, err := BlobListRequest(seedStr)
			if err != nil {
				fmt.Println(err)
				break
			}
			Do(r)
			Do(r)
		case 't':
			fmt.Println("Enter your seed:")
			seedStr, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}
			seedStr = seedStr[:len(seedStr)-2]
			r, err := CheckTimeConstraint(seedStr)
			if err != nil {
				fmt.Println(err)
				break
			}
			Do(r)
		}
	}
}
