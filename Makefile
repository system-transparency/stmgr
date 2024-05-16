default:
	true

# Keep things simple, no test driver script.
check:
	./tests/cert-keygen-test
	./tests/cert-ssh-test
	./tests/cert-openssl-test
	./tests/cert-agent-test
	./tests/hostconfig-check-test
	./tests/ospkg-create-test
	./tests/ospkg-sign-test
