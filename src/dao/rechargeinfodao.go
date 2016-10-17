package dao

type RechargeInfo struct {
	CompereId int64   `bson:"CompereId"`
	UserId int64    `bson:"UserId"`
	Sid int64    `bson:"Sid"`
	Ssid int64    `bson:"Ssid"`
	SignSid int64 `bson:"SignSid"`
	GiftId int64 `bson:"GiftId"`
	GiftValue int64 `bson:"GiftValue"`
	Platform int64 `bson:"Platform"`
	Timestamp int64 `bson:"Timestamp"`
	Days string `bson:"Days"`
	Hours string `bson:"Hour"`
	Group string `bson:"Group"`
}
