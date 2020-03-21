# Micro-conversion API
API for the microservice converting different formats to anothers.
For the moment, the only conversion possible is:
* image ==> image + piFF file (the piFF file is created and returned, it's not really a conversion)

## Home Link [/db]
Simple method to test whether the Go API is runing correctly

### [GET]
+ Response 200 (text/plain)  
	+ Body  
    	~~~
    	Welcome home!
    	~~~

## Create piFF from image [/convert/nothing]
This action returns a piFF file according to the image given in body.

### [POST]
This action has 2 negative responses defined:  
It will return a status 404 if a error occurs while converting the body to the appriopriate structure.  
It will return a status 500 if an error occurs in the Go service. This can happen in the image opening or in the file writing.  

+ Request (application/json)
	+ Body
		~~~
		{
			"Path": "path/to/image/in/server/image.png"
		}
		~~~

+ Response 200 (application/json)  
	+ Body  
		~~~
		{
            "Meta": {
                "Type": "line",
                "URL": ""
            },
            "Location": [
                {
                    "Type": "line",
                    "Polygon": [[0,0],[261,0],[261,343],[0,343]],
                    "Id": "loc_0"
                }
            ],
            "Data": [
                {
                    "Type": "line",
                    "LocationId": "loc_0",
                    "Value": "",
                    "Id": "0"
                }
            ],
            "Children": null,
            "Parent": 0
        }
        ~~~

+ Response 404 (text/plain)  
	+ Body  
		~~~
        [MICRO-CONVERSION] Couldn't read body
        ~~~

+ Response 500 (text/plain)  
	+ Body  
		~~~
        [MICRO-CONVERSION] {user-friendly error message}
        ~~~











