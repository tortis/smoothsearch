#include <iostream>
#include <NTL/GF2X.h>
#include <NTL/GF2XFactoring.h>

NTL_CLIENT

long getSmoothness(const GF2X& f) {
	static vec_pair_GF2X_long r;
	DDF(r, f);
	return r[r.length()-1].b;
}

int main(int argc, char* argv[])
{
	ZZ a, inc;
	if (argc == 3) {
		// Use commandline specified a and inc
		a = conv<ZZ>(argv[1]);
		inc = conv<ZZ>(argv[2]);
	} else {
		// Use default a and inc
		a = ZZ(113851762);
		inc = ZZ(1000000000);
	}

	// Modulus Polynomial
	GF2X m = GF2X(607, 1);
	SetCoeff(m, 105);
	SetCoeff(m, 0);

	// Base Polynomial 
	GF2X b = GF2X(606, 1);
	SetCoeff(b, 165);
	SetCoeff(b, 0);

	GF2X r;
	long s;
	long smoothest = 606;
	int iterations = 0;

	for (;;) {
		PowerMod(r, b, a, m);
		s = getSmoothness(r);
		if (s < smoothest) {
			smoothest = s;
			std::cout << "HIT:" << s << ":" << a << std::endl;
		}
		a += inc;
		iterations++;
		if (iterations == 10000) {
			iterations = 0;
			std::cout << "ITR:" << a << std::endl;
		}
	}
}
