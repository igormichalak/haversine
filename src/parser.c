#include <stdlib.h>
#include <stdio.h>
#include <inttypes.h>

#include "common.h"
#include "parser.h"

#define BUFFER_SIZE (16 * 1024)

typedef struct {
	FILE    *input;
	char    *buffer[BUFFER_SIZE];
	size_t   buffer_length;
	char     curr_char;
	char     peek_char;
} JSONParser;

