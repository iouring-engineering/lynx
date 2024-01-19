package main

const (
	INSERT_SHORT_LINK = `INSERT INTO LINKS (SHORT_CODE,DATA,WEB_URL,ANDROID,IOS,SOCIAL,EXPIRY) 
		VALUES (?,?,?,?,?,?,?)`
	GET_LINK_DATA = `SELECT DATA,WEB_URL,ANDROID,IOS,SOCIAL,EXPIRY FROM LINKS WHERE SHORT_CODE = ?`
)
