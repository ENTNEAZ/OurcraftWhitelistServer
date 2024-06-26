package WhitelistUtils

import (
	"WhiteListServer/Config"
	"encoding/json"
	"log"
	"os"
)

var WhitelistedUser []WhitelistProfileStruct

func ReadWhitelistFromJsonFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Read whitelist file failed,%v ,Ignore File", err)
		WhitelistedUser = make([]WhitelistProfileStruct, 0)
		return
	}
	err = json.Unmarshal(content, &WhitelistedUser)
	if err != nil {
		log.Printf("Unmarshal whitelist file failed,%v", err)
		panic(err)
	}

	log.Printf("Read whitelist file success,File Path:%s", filePath)
}

func SaveWhitelistToJsonFile(filePath string) {
	content, err := json.Marshal(WhitelistedUser)
	if err != nil {
		log.Printf("Marshal whitelist file failed,%v", err)
		panic(err)
	}
	err = os.WriteFile(filePath, content, 0666)
	if err != nil {
		log.Printf("Write whitelist file failed,%v", err)
		panic(err)
	}
	log.Printf("Save whitelist file success,File Path:%s", filePath)

}

func CheckIfInWhitelist(name string) bool {
	for _, v := range WhitelistedUser {
		if v.Name == name {
			log.Printf("Lookup Warning: User:'%s' already exists", name)
			return true
		}
	}
	return false
}

func AddToWhitelist(name, uuid, contactMethod, contactID string) {
	WhitelistedUser = append(WhitelistedUser, WhitelistProfileStruct{
		UUID:          uuid,
		Name:          name,
		ContactMethod: contactMethod,
		ContactID:     contactID,
	})
	log.Printf("Add user '%s' to whitelist", name)
	SaveWhitelistToJsonFile(Config.GetConfig().WhitelistFilePath)
}

func CheckIfQQInWhitelist(qq string) bool {
	for _, v := range WhitelistedUser {
		if v.ContactMethod == "qq" && v.ContactID == qq {
			log.Printf("Lookup Warning: QQ:'%s' already exists", qq)
			return true
		}
	}
	return false
}
