package fbitsdk

import (
	//"bytes"
	//"encoding/json"
	//"crypto"
	//"io"
	//"net/http"
	"net/url"
	//"strings"
	"time"
)

func sdkBodyFat(form url.Values, uid string) (string, error) {
	urlString := srcUrl1
	urlString += uid
	urlString += heartRate
	nowStr := time.Now().String()
	urlString += nowStr[0:10]
	urlString += ".json"
	return urlString, nil
}

// Body Weight
//There are three acceptable formats for retrieving body fat log data:

//GET https://api.fitbit.com/1/user/[user-id]/body/log/fat/date/[date].json
//GET https://api.fitbit.com/1/user/[user-id]/body/log/fat/date/[date]/[period].json
//GET https://api.fitbit.com/1/user/[user-id]/body/log/fat/date/[base-date]/[end-date].json
//user-id	The encoded ID of the user. Use "-" (dash) for current logged-in user.
//date	The date in the format yyyy-MM-dd.
//base-date	The end date when period is provided; range start date when a date range is provided. In the format yyyy-MM-dd or today.
//period	The date range period. One 1d, 7d, 1w, 1m.
//end-date	Range end date when date range is provided. Note: The range should not be longer than 31 days.
//Request Headers
//Accept-Language	optional	The measurement unit system to use for response values.
//Response Fields
//date	Log entry date; in the format yyyy-MM-dd.
//fat	Body fat percentage; in the format X.XX.
//logId	Body Fat Log IDs are unique to the user, but not globally unique.
//time	Time of the measurement; hours and minutes in the format HH:mm:ss, set to the last second of the day if not provided.
//source	The source of the fat log; the field is optional.
//Example Response
/*
{
    "fat":[
        {
            "date":"2012-03-05",
            "fat":14,
            "logId":1330991999000,
            "time":"23:59:59",
            "source": "API"
        },
        {
            "date":"2012-03-05",
            "fat":13.5,
            "logId":1330991999000,
            "time":"21:20:59",
            "source":"Aria"
        }
    ]
}
*/
