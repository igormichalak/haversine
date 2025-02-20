#ifndef COMMON_H
#define COMMON_H

#define PI 3.14159265358979323846264338327950288419716939937510582097494459
#define EARTH_RADIUS 6372.8
#define RADIANS_PER_DEGREE (PI / 180.0)

typedef double f64;

typedef struct {
	f64 x0;
	f64 y0;
	f64 x1;
	f64 y1;
} PointPair;

#endif
