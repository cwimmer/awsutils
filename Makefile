pre-commit:
	pre-commit autoupdate
	pre-commit run --all-files

bump:
	cz bump
	git push --tags
