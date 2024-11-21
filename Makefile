

.PHONY: release
release:
	dagger call release \
		--one-password-service-account-production env:OP_SERVICE_ACCOUNT_PRODUCTION \
		--version $(version) \
		--github-token env:GITHUB_TOKEN \
		--progress plain
