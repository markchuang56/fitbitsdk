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

func sdkSleep(form url.Values, uid string) (string, error) {
	urlString := srcUrl2
	urlString += uid
	urlString += heartRate
	nowStr := time.Now().String()
	urlString += nowStr[0:10]
	urlString += ".json"
	return urlString, nil
}

// SLEEP

//Resource URL
//GET https://api.fitbit.com/1.2/user/[user-id]/sleep/date/[date].json
//URL parameters:

//user-id	The ID of the user. Use "-" (dash) for current logged-in user.
//date	The date of records to be returned. In the format yyyy-MM-dd.
//Example Request
//GET https://api.fitbit.com/1.2/user/-/sleep/date/2017-04-02.json
//Example Response
//Note: The text within the brackets <> is a descriptive placeholder for a value or repeated elements.
/*
{
    "sleep": [
        {
            "dateOfSleep": "2017-04-02",
            "duration": <value in milliseconds>,
            "efficiency": <value>,
            "isMainSleep": true,
            "levels": {
                "summary": {
                    "deep": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    },
                    "light": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    },
                    "rem": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    },
                    "wake": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    }
                },
                "data": [
                    {
                        "datetime": "2017-04-01T23:58:30.000",
                        "level": "wake",
                        "seconds": <value>
                    },
                    {
                        "datetime": "2017-04-02T00:16:30.000",
                        "level": "rem",
                        "seconds": <value>
                    },
                    <...>
                ],
                "shortData": [
                    {
                        "datetime": "2017-04-02T05:58:30.000",
                        "level": "wake",
                        "seconds": <value>
                    },
                    <...>
                ]
            },
            "logId": <value>,
            "minutesAfterWakeup": <value>,
            "minutesAsleep": <value>,
            "minutesAwake": <value>,
            "minutesToFallAsleep": <value>, // this is generally 0 for autosleep created sleep logs
            "startTime": "2017-04-01T23:58:30.000",
            "timeInBed": <value in minutes>,
            "type": "stages"
        },
        {
            "dateOfSleep": "2017-04-02",
            "duration": <value in milliseconds>,
            "efficiency": <value>,
            "isMainSleep": false,
            "levels": {
                "data": [
                    {
                        "dateTime": "2017-04-02T12:06:00.000",
                        "level": "asleep",
                        "seconds": <value>
                    },
                    {
                        "dateTime": "2017-04-02T12:13:00.000",
                        "level": "restless",
                        "seconds": <value>
                    },
                    {
                        "dateTime": "2017-04-02T12:14:00.000",
                        "level": "awake",
                        "seconds": <value>
                    },
                    <...>
                ],
                "summary": {
                    "asleep": {
                        "count": 0, // this field should not be used for "asleep" summary info
                        "minutes": <value>
                    },
                    "awake": {
                        "count": <value>,
                        "minutes": <value>
                    },
                    "restless": {
                        "count": <value>,
                        "minutes": <value>
                    }
                }
            },
            "logId": <value>,
            "minutesAfterWakeup": <value>,
            "minutesAsleep": <value>,
            "minutesAwake": <value>,
            "minutesToFallAsleep": <value>, // this is generally 0 for autosleep created sleep logs
            "startTime": "2017-04-02T12:06:00.000",
            "timeInBed": <value in minutes>,
            "type": "classic"
        }
    ],
    "summary": {
        "totalMinutesAsleep": <value>,
        "totalSleepRecords": 2,
        "totalTimeInBed": <value in minutes>
    }
}

Note: Some processing is asynchronous. If the system is still processing one or more sleep logs that should be in the response when this API is queried, the response will indicate a retry duration given in milliseconds. The "meta" response may evolve with additional fields in the future; API clients should be prepared to ignore any new object properties they do not recognize.

{
    "meta": {
        "retryDuration": 3000,
        "state": "pending"
    }
}
*/

// Get Sleep Logs by Date Range

//esource URL
//GET https://api.fitbit.com/1.2/user/[user-id]/sleep/date/[startDate]/[endDate].json
//URL parameters:

//user-id	The ID of the user. Use "-" (dash) for current logged-in user.
//startDate	The date for the first sleep log to be returned. In the format yyyy-MM-dd. This date is inclusive
//endDate	The date for the end sleep log to be returned. In the format yyyy-MM-dd. This date is inclusive.
//Example Request
//GET https://api.fitbit.com/1.2/user/-/sleep/date/2017-04-02/2017-04-08.json
//Example Response
//Note: The text within the brackets <> is a descriptive placeholder for a value or repeated elements.
/*
{
    "sleep": [
        {
            "dateOfSleep": "2017-04-02",
            "duration": <value in milliseconds>,
            "efficiency": <value>,
            "isMainSleep": <true|false>,
            "levels": {
                "summary": {
                    "deep": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    },
                    "light": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    },
                    "rem": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    },
                    "wake": {
                        "count": <value>,
                        "minutes": <value>,
                        "thirtyDayAvgMinutes": <value>
                    }
                },
                "data": [
                    {
                        "datetime": "2017-04-01T23:58:30.000",
                        "level": "wake",
                        "seconds": <value>
                    },
                    {
                        "datetime": "2017-04-02T00:16:30.000",
                        "level": "light",
                        "seconds": <value>
                    },
                    <...>
                ],
                "shortData": [
                    {
                        "datetime": "2017-04-02T05:58:30.000",
                        "level": "wake",
                        "seconds": <value>
                    },
                    <...>
                ]
            },
            "logId": <value>,
            "minutesAfterWakeup": <value>,
            "minutesAsleep": <value>,
            "minutesAwake": <value>,
            "minutesToFallAsleep": <value>, // this is generally 0 for autosleep created sleep logs
            "startTime": "2017-04-01T23:58:30.000",
            "timeInBed": <value in minutes>,
            "type": "stages"
        },
        <...>
    ]
}
Note: Some processing is asynchronous. If the system is still processing one or more sleep logs that should be in the response when this API is queried, the response will indicate a retry duration given in milliseconds. The "meta" response may evolve with additional fields in the future; API clients should be prepared to ignore any new object properties they do not recognize.

{
    "meta": {
        "retryDuration": 3000,
        "state": "pending"
    }
}
*/
