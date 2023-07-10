package cloudranges

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// command to use to check tls records:	tlsx -san -cn -silent

// https://kaeferjaeger.gay/?dir=ip-ranges

// amazon 		✅
// bing 		✅
// digitalocean ✅
// facebook
// github 		✅
// google 		✅
// linode 		✅
// microsoft
// oracle 		✅
// telegram 	✅
// twitter

// TODO: make them to go function (concurrency)

type bingbotResponse struct {
	CreationTime string `json:"creationTime"`
	Prefixes     []struct {
		Ipv4Prefix string `json:"ipv4Prefix"`
	} `json:"prefixes"`
}

func BingbotRanges() ([]string, error) {
	resp, err := http.Get("https://www.bing.com/toolbox/bingbot.json")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// defer resp.Body.Close()

	var ips bingbotResponse

	if err := json.Unmarshal([]byte(b), &ips); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var ipList []string
	for _, v := range ips.Prefixes {
		// fmt.Println(v.Ipv4Prefix)
		ipList = append(ipList, v.Ipv4Prefix)
	}

	return ipList, nil
}

type googleCloudResponse struct {
	SyncToken    string `json:"syncToken"`
	CreationTime string `json:"creationTime"`
	Prefixes     []struct {
		Ipv4Prefix string `json:"ipv4Prefix,omitempty"`
		Service    string `json:"service"`
		Scope      string `json:"scope"`
		Ipv6Prefix string `json:"ipv6Prefix,omitempty"`
	} `json:"prefixes"`
}

func GoogleCloudRanges(s string) ([]string, error) {
	// IP ranges that Google makes available to users on the internet
	resp1, err1 := http.Get("https://www.gstatic.com/ipranges/cloud.json")
	if err1 != nil {
		return nil, fmt.Errorf("%v", err1)
	}

	defer resp1.Body.Close()
	// Global and regional external IP address ranges for customers' Google Cloud resources
	resp2, err2 := http.Get("https://www.gstatic.com/ipranges/goog.json")
	if err2 != nil {
		return nil, fmt.Errorf("%v", err2)
	}
	defer resp2.Body.Close()

	var ips googleCloudResponse

	b1, err := io.ReadAll(resp1.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal([]byte(b1), &ips); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	b2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal([]byte(b2), &ips); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var ipsList []string

	for _, val := range ips.Prefixes {
		if s == "v4" {
			ipsList = append(ipsList, val.Ipv4Prefix)
		} else if s == "v6" {
			if val.Ipv6Prefix != "" {
				ipsList = append(ipsList, val.Ipv6Prefix)
			}
		} else if s == "all" {
			if val.Ipv6Prefix != "" {
				ipsList = append(ipsList, val.Ipv6Prefix)
			}
			ipsList = append(ipsList, val.Ipv4Prefix)
		}
	}

	return ipsList, nil
}

func GoogleBotRanges(s string) ([]string, error) {

	resp, err := http.Get("https://developers.google.com/search/apis/ipranges/googlebot.json")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer resp.Body.Close()
	var ips googleCloudResponse

	b1, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal([]byte(b1), &ips); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var ipsList []string

	for _, val := range ips.Prefixes {
		if s == "v4" {
			ipsList = append(ipsList, val.Ipv4Prefix)
		} else if s == "v6" {
			if val.Ipv6Prefix != "" {
				ipsList = append(ipsList, val.Ipv6Prefix)
			}
		} else if s == "all" {
			if val.Ipv6Prefix != "" {
				ipsList = append(ipsList, val.Ipv6Prefix)
			}
			if val.Ipv4Prefix != "" {
				ipsList = append(ipsList, val.Ipv4Prefix)
			}
		}
	}
	return ipsList, nil
}

type awsCloudResponse struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix           string `json:"ip_prefix"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
	} `json:"prefixes"`
	Ipv6Prefixes []struct {
		Ipv6Prefix         string `json:"ipv6_prefix"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
	} `json:"ipv6_prefixes"`
}

func AwsRanges(s string) ([]string, error) {
	resp, err := http.Get("https://ip-ranges.amazonaws.com/ip-ranges.json")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer resp.Body.Close()

	var ips awsCloudResponse

	b1, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal([]byte(b1), &ips); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var ipsList []string

	if s == "all" {
		for _, val := range ips.Ipv6Prefixes {
			ipsList = append(ipsList, val.Ipv6Prefix)
		}

		for _, val := range ips.Prefixes {
			ipsList = append(ipsList, val.IPPrefix)
		}
	} else if s == "v4" {
		for _, val := range ips.Prefixes {
			ipsList = append(ipsList, val.IPPrefix)
		}
	} else if s == "v6" {
		for _, val := range ips.Ipv6Prefixes {
			ipsList = append(ipsList, val.Ipv6Prefix)
		}
	}

	return ipsList, nil
}

type oracleCloudResponse struct {
	LastUpdatedTimestamp string `json:"last_updated_timestamp"`
	Regions              []struct {
		Region string `json:"region"`
		Cidrs  []struct {
			Cidr string   `json:"cidr"`
			Tags []string `json:"tags"`
		} `json:"cidrs"`
	} `json:"regions"`
}

func OracleRanges() ([]string, error) {
	resp, err := http.Get("https://docs.oracle.com/en-us/iaas/tools/public_ip_ranges.json")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer resp.Body.Close()
	var ips oracleCloudResponse

	b1, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal([]byte(b1), &ips); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var ipsList []string

	for _, region := range ips.Regions {
		for _, cidrs := range region.Cidrs {
			fmt.Println(cidrs.Cidr)
			ipsList = append(ipsList, cidrs.Cidr)
		}
	}

	return ipsList, nil
}

func getDigitalOceanRanges() (string, error) {
	filePath := "/tmp/digitalOceanRanges.txt"

	out, ferr := os.Create(filePath)
	if ferr != nil {
		return "", fmt.Errorf("%v", ferr)
	}

	defer out.Close()

	resp, rerr := http.Get("https://www.digitalocean.com/geo/google.csv")
	if rerr != nil {
		return "", fmt.Errorf("%v", rerr)
	}

	defer resp.Body.Close()

	_, cerr := io.Copy(out, resp.Body)
	if cerr != nil {
		return "", fmt.Errorf("%v", cerr)
	}

	return filePath, nil
}

func DigitalOceanRanges(s string) ([]string, error) {
	fp, fperr := getDigitalOceanRanges()
	if fperr != nil {
		// fmt.Println("Error opening CSV file:", err)
		return nil, fmt.Errorf("%v", fperr)
	}

	file, foerr := os.Open(fp)
	if foerr != nil {
		// fmt.Println("Error opening CSV file:", err)
		return nil, fmt.Errorf("%v", foerr)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	records, raerr := reader.ReadAll()
	if raerr != nil {
		// fmt.Println("Error reading CSV file:", err)
		return nil, fmt.Errorf("%v", raerr)
	}

	var ipList []string

	for _, record := range records {
		// Access each field of the record
		if s == "all" {
			if strings.Contains(record[0], ".") {
				ipList = append(ipList, record[0])
			}
			if strings.Contains(record[0], ":") {
				ipList = append(ipList, record[0])
			}
		} else if s == "v4" {
			if strings.Contains(record[0], ".") {
				ipList = append(ipList, record[0])
			}
		} else if s == "v6" {
			if strings.Contains(record[0], ":") {
				ipList = append(ipList, record[0])
			}
		}

	}
	return ipList, nil
}

func LinodeRanges(s string) ([]string, error) {
	resp, err := http.Get("https://geoip.linode.com")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.Comment = '#'

	records, raerr := reader.ReadAll()
	if raerr != nil {
		// fmt.Println("Error reading CSV file:", err)
		return nil, fmt.Errorf("%v", raerr)
	}

	var ipList []string

	for _, record := range records {
		// Access each field of the record
		if s == "v4" {
			if strings.Contains(record[0], ".") {
				// fmt.Println(record[0])
				ipList = append(ipList, record[0])
			}
		} else if s == "v6" {
			if strings.Contains(record[0], ":") {
				// fmt.Println(record[0])
				ipList = append(ipList, record[0])
			}
		} else if s == "all" {
			ipList = append(ipList, record[0])
		}

	}
	return ipList, nil
}

func TelegramRanges(s string) ([]string, error) {
	// https://core.telegram.org/resources/cidr.txt
	resp, err := http.Get("https://core.telegram.org/resources/cidr.txt")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var ipList []string

	// Read each line of the response
	for scanner.Scan() {
		line := scanner.Text()
		if s == "v4" {
			if strings.Contains(line, ".") {
				ipList = append(ipList, line)
			}
		} else if s == "v6" {
			if strings.Contains(line, ":") {
				ipList = append(ipList, line)
			}
		} else if s == "all" {
			ipList = append(ipList, line)
		}

	}

	if serr := scanner.Err(); err != nil {
		fmt.Println("Error reading response:", serr)
		return nil, serr
	}
	return ipList, nil
}

type githubResponse struct {
	VerifiablePasswordAuthentication bool `json:"verifiable_password_authentication"`
	SSHKeyFingerprints               struct {
		Sha256Ecdsa   string `json:"SHA256_ECDSA"`
		Sha256Ed25519 string `json:"SHA256_ED25519"`
		Sha256Rsa     string `json:"SHA256_RSA"`
	} `json:"ssh_key_fingerprints"`
	SSHKeys    []string `json:"ssh_keys"`
	Hooks      []string `json:"hooks"`
	Web        []string `json:"web"`
	API        []string `json:"api"`
	Git        []string `json:"git"`
	Packages   []string `json:"packages"`
	Pages      []string `json:"pages"`
	Importer   []string `json:"importer"`
	Actions    []string `json:"actions"`
	Dependabot []string `json:"dependabot"`
	Domains    struct {
		Website    []string `json:"website"`
		Codespaces []string `json:"codespaces"`
		Copilot    []string `json:"copilot"`
		Packages   []string `json:"packages"`
	} `json:"domains"`
}

func GithubRanges(s string) ([]string, error) {
	resp, err := http.Get("https://api.github.com/meta")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	//defer resp.Body.Close()

	var ips githubResponse

	b1, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	if err := json.Unmarshal([]byte(b1), &ips); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	combinded := append(ips.Hooks, ips.Web...)
	combinded = append(combinded, ips.API...)
	combinded = append(combinded, ips.Git...)
	combinded = append(combinded, ips.Packages...)
	combinded = append(combinded, ips.Pages...)
	combinded = append(combinded, ips.Importer...)
	combinded = append(combinded, ips.Actions...)
	combinded = append(combinded, ips.Dependabot...)

	var ipLists []string

	for _, v := range combinded {
		if s == "v4" {
			if strings.Contains(v, ".") {
				ipLists = append(ipLists, v)
			}
		} else if s == "v6" {
			if strings.Contains(v, ":") {
				ipLists = append(ipLists, v)
			}
		} else if s == "all" {
			ipLists = append(ipLists, v)
		}
	}

	return ipLists, nil
}
