package handler

import (
	"WhiteListServer/Config"
	"WhiteListServer/MojangUtils"
	"WhiteListServer/QQGroupUtils"
	"WhiteListServer/WhitelistUtils"
	"log"
	"net/http"
)

func getRealIP(r *http.Request) string {
	if r.Header.Get("X-Real-IP") != "" {
		return r.Header.Get("X-Real-IP")
	}
	return r.RemoteAddr
}

func ApplyWhitelist(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	contactMethod := r.URL.Query().Get("contactmethod")
	contactID := r.URL.Query().Get("contactid")
	ip := getRealIP(r)

	// 中文模式下不会发送信息
	//其他奇怪的QQ contactMethod
	if contactMethod == "" || contactMethod == "QQ" || contactMethod == "Qq" || contactMethod == "qQ" {
		contactMethod = "qq"
	}

	//确认不在白名单钟
	switch {
	case contactMethod == "":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"请填写联系方式\",\"success\":false}"))
		log.Printf("[IP:%s]User %s Applied,but contactMethod is empty", ip, name)
		return
	case contactID == "":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"请填写联系方式\",\"success\":false}"))
		log.Printf("[IP:%s]User %s Applied,but contactID is empty", ip, name)
		return
	case name == "":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\":\"请填写名字\",\"success\":false}"))
		log.Printf("[IP:%s]User %s Applied,but name is empty", ip)
		return
	}
	log.Printf("[IP:%s]User %s Applied,name:\"%s\",contactMethod:\"%s\",contactID:\"%s\"", ip, name, name, contactMethod, contactID)

	//先确认是否在白名单中
	if WhitelistUtils.CheckIfInWhitelist(name) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("{\"error\":\"你已经在白名单申请列表中了哦，请不要重复申请\",\"success\":false}"))
		log.Printf("[IP:%s]User %s Applied,but already in the apply list", ip, name)
		return
	}

	//若method是qq 则检查是否已经有人用这个qq申请过了
	if contactMethod == "qq" {
		if WhitelistUtils.CheckIfQQInWhitelist(contactID) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("{\"error\":\"这个QQ已经有人申请过了哦，请不要重复申请\",\"success\":false}"))
			log.Printf("[IP:%s]User %s Applied,but QQ %s already in the apply list", ip, name, contactID)
			return
		}

	}

	//若method是qq 检查这个人是不是在q群里
	if contactMethod == "qq" {
		i, _, err := QQGroupUtils.CheckIfUserInGroup(contactID)
		if !i {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("{\"error\":\"你不在QQ群中，请先加QQ群:" + Config.GetConfig().GroupID + "后再申请哦，如果有疑问请联系群内管理\",\"success\":false}"))
			log.Printf("[IP:%s]User %s Applied,but QQ %s not in the QQ group", ip, name, contactID)
			return
		}
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\":\"QQ群验证失败，请联系管理员\",\"success\":false}"))
			log.Printf("[IP:%s]User %s Applied,but QQ %s check failed", ip, name, contactID)
			return
		}
	}

	//获取用户相关信息
	profile, err := MojangUtils.GetProfileByUserName(name)
	if err != nil {
		switch {
		case err.Error() == "not found":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{\"error\":\"无法找到此ID玩家,请确认你的输入\",\"success\":false}"))
			log.Printf("[IP:%s]User %s Applied,but user not found(via MojangUtils)", ip, name)
			return
		default:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"error\":\"Mojang服务器错误\",\"success\":false}"))
			log.Printf("[IP:%s]User %s Applied,but MojangUtils error:%s", ip, name, err)
			return
		}
	}
	log.Printf("[IP:%s]User %s ,UUID:%s", ip, name, profile.UUID)
	WhitelistUtils.AddToWhitelist(name, profile.UUID, contactMethod, contactID)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte("{\"success\":true}"))
	return

}
