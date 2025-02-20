#include <math.h>

#include "common.h"
#include "haversine.h"

static f64 square(f64 n) {
	return n * n;
}

f64 reference_haversine(f64 x0, f64 y0, f64 x1, f64 y1, f64 radius)
{
	f64 lat1 = y0;
	f64 lat2 = y1;
	f64 lon1 = x0;
	f64 lon2 = x1;

	f64 d_lat = (lat2 - lat1) * RADIANS_PER_DEGREE;
	f64 d_lon = (lon2 - lon1) * RADIANS_PER_DEGREE;

	lat1 *= RADIANS_PER_DEGREE;
	lat2 *= RADIANS_PER_DEGREE;

	f64 a = square(sin(d_lat / 2.0)) + cos(lat1) * cos(lat2) * square(sin(d_lon / 2.0));
	f64 c = 2.0 * asin(sqrt(a));

	return c * radius;
}
