#include <stdlib.h>
#include <stdio.h>

#include "common.h"
#include "parser.h"

int main(int argc, char **argv) {
	if (argc < 2) {
		fputs("No filename!\n", stderr);
		return EXIT_FAILURE;
	}
	FILE *fp = fopen(argv[1], "r");
	if (fp == NULL) {
		fprintf(stderr, "Cannot open %s!", argv[1]);
		return EXIT_FAILURE;
	}

	fclose(fp);

	return EXIT_SUCCESS;
}
