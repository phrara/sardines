package config

import (
	"sardines/tool"
	"encoding/json"
	"errors"
	"log"
	mrand "math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p-core/crypto"
)

var cpath string
var kpath string
var Ktab string
var Dir string

func init() {
	wd, _ := os.Getwd()
	Dir = filepath.Join(wd, "/values")
	kpath = filepath.Join(Dir, "/priv_key")
	Ktab = filepath.Join(Dir, "/key_tab.db")
	cpath = filepath.Join(Dir, "/config.json")
}

type Config struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	RandomSeed int64  `json:"randomSeed"`
	// Private Key
	PrvKey crypto.PrivKey `json:"prvKey"`
	// Bootstrap Node
	BootstrapNode string `json:"bootstrapNode"`
}

func (c *Config) Save() bool {
	c.PrvKey = nil
	b, _ := json.Marshal(*c)
	err := tool.WriteFile(b, cpath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func (c *Config) LoadAll() *Config {
	data, err := tool.LoadFile(cpath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	pk, _ := tool.LoadFile(kpath)
	if len(pk) == 0 {
		// prvKey is empty
		log.Printf("you haven't configure your own private key, so that your Peer ID could be variable")
		c.PrvKey = nil
	} else {
		c.PrvKey, _ = crypto.UnmarshalPrivateKey(pk)
	}
	return c
}

func (c *Config) Load() *Config {
	data, err := tool.LoadFile(cpath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	c.PrvKey = nil
	return c
}

func (c *Config) AddrString() string {
	// "/ip4/127.0.0.1/tcp/2000"
	return "/ip4/" + c.IP + "/tcp/" + c.Port
}

// 获取私钥

func (c *Config) GenKey() error {
	if c.RandomSeed > 0 {
		r := mrand.New(mrand.NewSource(c.RandomSeed))
		prvKey, _, err := crypto.GenerateRSAKeyPair(2048, r)
		if err != nil {
			return err
		}
		c.PrvKey = prvKey
		return nil
	} else {
		return errors.New("the RandomSeed is not a positive number")
	}
}

func (c *Config) SaveKey() error {
	if c.PrvKey != nil {
		b, err := crypto.MarshalPrivateKey(c.PrvKey)
		if err != nil {
			return err
		}
		if err = tool.WriteFile(b, kpath); err != nil {
			return err
		}
	}
	return nil
}

// New Get a Configuration
func New(username, pwd, ipAddr string, rs int64, bn string) (*Config, error) {
	strs := strings.Split(ipAddr, ":")
	if strs[0] == "" {
		strs[0] = "127.0.0.1"
	}

	if b := ipFormatCheck(strs); b {
		return &Config{
			Username:      username,
			Password:      pwd,
			IP:            strs[0],
			Port:          strs[1],
			RandomSeed:    rs,
			PrvKey:        nil,
			BootstrapNode: bn,
		}, nil
	} else {
		return nil, errors.New("the format of ipAddr is wrong")
	}
}
 
func ipFormatCheck(ipAddr []string) bool {
	compile, _ := regexp.Compile(`((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)`)
	if b := compile.MatchString(ipAddr[0]); b {
		port, err := strconv.ParseInt(ipAddr[1], 10, 64)
		if err != nil || port < 0 || port > 65535 {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}
