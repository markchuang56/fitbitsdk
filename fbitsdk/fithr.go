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

func sdkHeartRate(form url.Values, uid string) (string, error) {
	urlString := srcUrl1
	urlString += uid
	urlString += heartRate
	nowStr := time.Now().String()
	urlString += nowStr[0:10]
	urlString += ".json"
	return urlString, nil
}

// Heart Rate
//GET https://api.fitbit.com/1/user/[user-id]/activities/heart/date/[date]/[period].json
//GET https://api.fitbit.com/1/user/[user-id]/activities/heart/date/[base-date]/[end-date].json

//Example Request
//GET https://api.fitbit.com/1/user/-/activities/heart/date/today/1d.json
//Example Response
/*
{
    "activities-heart": [
        {
            "dateTime": "2015-08-04",
            "value": {
                "customHeartRateZones": [],
                "heartRateZones": [
                    {
                        "caloriesOut": 740.15264,
                        "max": 94,
                        "min": 30,
                        "minutes": 593,
                        "name": "Out of Range"
                    },
                    {
                        "caloriesOut": 249.66204,
                        "max": 132,
                        "min": 94,
                        "minutes": 46,
                        "name": "Fat Burn"
                    },
                    {
                        "caloriesOut": 0,
                        "max": 160,
                        "min": 132,
                        "minutes": 0,
                        "name": "Cardio"
                    },
                    {
                        "caloriesOut": 0,
                        "max": 220,
                        "min": 160,
                        "minutes": 0,
                        "name": "Peak"
                    }
                ],
                "restingHeartRate": 68
            }
        }
    ]
}
*/

//user-id	The encoded ID of the user. Use "-" (dash) for current logged-in user.
//base-date	The range start date, in the format yyyy-MM-dd or today.
//end-date	The end date of the range.
//date	The end date of the period specified in the format yyyy-MM-dd or today.
//period	The range for which data will be returned. Options are 1d, 7d, 30d, 1w, 1m.

//GET https://api.fitbit.com/1/user/-/activities/heart/date/[date]/[end-date]/[detail-level].json
//GET https://api.fitbit.com/1/user/-/activities/heart/date/[date]/[end-date]/[detail-level]/time/[start-time]/[end-time].json
//GET https://api.fitbit.com/1/user/-/activities/heart/date/[date]/1d/[detail-level].json`
//GET https://api.fitbit.com/1/user/-/activities/heart/date/[date]/1d/[detail-level]/time/[start-time]/[end-time].json
//date	The date, in the format yyyy-MM-dd or today.
//detail-level	Number of data points to include. Either 1sec or 1min. Optional.
//start-time	The start of the period, in the format HH:mm. Optional.
//end-time	The end of the period, in the format HH:mm. Optional.
//Example Request
//GET https://api.fitbit.com/1/user/-/activities/heart/date/today/1d/1sec/time/00:00/0

//Example Response
/*
{
    "activities-heart": [
        {
            "customHeartRateZones": [],
            "dateTime": "today",
            "heartRateZones": [
                {
                    "caloriesOut": 2.3246,
                    "max": 94,
                    "min": 30,
                    "minutes": 2,
                    "name": "Out of Range"
                },
                {
                    "caloriesOut": 0,
                    "max": 132,
                    "min": 94,
                    "minutes": 0,
                    "name": "Fat Burn"
                },
                {
                    "caloriesOut": 0,
                    "max": 160,
                    "min": 132,
                    "minutes": 0,
                    "name": "Cardio"
                },
                {
                    "caloriesOut": 0,
                    "max": 220,
                    "min": 160,
                    "minutes": 0,
                    "name": "Peak"
                }
            ],
            "value": "64.2"
        }
    ],
    "activities-heart-intraday": {
        "dataset": [
            {
                "time": "00:00:00",
                "value": 64
            },
            {
                "time": "00:00:10",
                "value": 63
            },
            {
                "time": "00:00:20",
                "value": 64
            },
            {
                "time": "00:00:30",
                "value": 65
            },
            {
                "time": "00:00:45",
                "value": 65
            }
        ],
        "datasetInterval": 1,
        "datasetType": "second"
    }
}
*/
