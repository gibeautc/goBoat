package boat

import "math"


/*
lat,lon are decimal degrees  (float64)
distance is measured in inches (int64)
angles are measred in degrees (int)
speed in mph
*/

func distanceBetween(lat1 float64, lon1 float64, lat2 float64, lon2 float64) (int64, int) {
	//todo:write tests and test
	/*
	   Distance JavaScript
	   var R = 6371e3; // metres
	   var φ1 = lat1.toRadians();
	   var φ2 = lat2.toRadians();
	   var Δφ = (lat2-lat1).toRadians();
	   var Δλ = (lon2-lon1).toRadians();

	   var a = Math.sin(Δφ/2) * Math.sin(Δφ/2) +
	           Math.cos(φ1) * Math.cos(φ2) *
	           Math.sin(Δλ/2) * Math.sin(Δλ/2);
	   var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));

	   var d = R * c;
	*/

	/*
	 Bearing
	 where 	φ1,λ1 is the start point, φ2,λ2 the end point (Δλ is the difference in longitude)
	 var y = Math.sin(λ2-λ1) * Math.cos(φ2);
	var x = Math.cos(φ1)*Math.sin(φ2) -
	        Math.sin(φ1)*Math.cos(φ2)*Math.cos(λ2-λ1);
	var brng = Math.atan2(y, x).toDegrees();
	*/

	R := 6371000.0 //m
	latDelta := toRadians(lat2 - lat1)
	lonDelta := toRadians(lon2 - lon1)
	lat1 = toRadians(lat1)
	lat2 = toRadians(lat2)
	a := math.Sin(latDelta/2)*math.Sin(latDelta/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(lonDelta/2)*math.Sin(lonDelta/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c      // in m
	d = 39.3701 * d //inches

	y := math.Sin(lonDelta) * math.Cos(lat2)
	x := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(lonDelta)
	brng := math.Atan2(y, x)

	//return distance, bearing
	return int64(d), int(toDegrees(brng))
}

func toRadians(deg float64) float64 {
	return deg / 360 * 2 * math.Pi
}

func toDegrees(rad float64) float64 {
	return rad / (2 * math.Pi) * 360
}

func getCords(lat float64, lon float64, distance int64, direction int) (float64, float64) {
	//todo:write tests and test
	/*
	   python example
	   import math

	   R = 6378.1 #Radius of the Earth
	   brng = 1.57 #Bearing is 90 degrees converted to radians.
	   d = 15 #Distance in km

	   #lat2  52.20444 - the lat result I'm hoping for
	   #lon2  0.36056 - the long result I'm hoping for.

	   lat1 = math.radians(52.20472) #Current lat point converted to radians
	   lon1 = math.radians(0.14056) #Current long point converted to radians

	   lat2 = math.asin( math.sin(lat1)*math.cos(d/R) +
	        math.cos(lat1)*math.sin(d/R)*math.cos(brng))

	   lon2 = lon1 + math.atan2(math.sin(brng)*math.sin(d/R)*math.cos(lat1),
	                math.cos(d/R)-math.sin(lat1)*math.sin(lat2))

	   lat2 = math.degrees(lat2)
	   lon2 = math.degrees(lon2)

	   print(lat2)
	   print(lon2)
	*/
	R := 6378.1 //km
	angle := toRadians(float64(direction))
	d := float64(distance) * .0000254
	startLat := toRadians(lat)
	startLon := toRadians(lon)
	lat2 := math.Asin(math.Sin(startLat)*math.Cos(d/R) + math.Cos(startLat)*math.Sin(d/R) + math.Cos(angle))
	lon2 := startLon + math.Atan2(math.Sin(angle)*math.Sin(d/R)*math.Cos(startLat)*math.Cos(d/R), math.Cos(d/R)-math.Sin(startLat)*math.Sin(lat2))
	return toDegrees(lat2), toDegrees(lon2)
}


func sliceContainsUint32(lst []uint32,obj uint32) bool{
	for x:=0;x<len(lst);x++{
		if lst[x]==obj{
			return true
		}
	}
	return false
}