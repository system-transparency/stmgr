default:
	true

# Keep things simple, no test driver script.
check:
	./tests/keygen-cert-test
	./tests/ssh-cert-test
	./tests/openssl-cert-test
	./tests/agent-cert-test
	./tests/hostconfig-check-test
	./tests/ospkg-sign-test
