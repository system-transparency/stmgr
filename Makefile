default:
	true

# Keep things simple, no test driver script.
check:
	./tests/keygen-cert-test
	./tests/hostconfig-check-test
