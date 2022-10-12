all: dockerBuild dockerRun

dockerBuild:
	sh ./scripts/build.sh

dockerRun:
	sh ./scripts/run.sh

release:
	bash ./scripts/release.sh $(version)
