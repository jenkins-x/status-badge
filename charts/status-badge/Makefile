CHART_REPO := https://chartmuseum.jx.charts-demo.cloudbees.com
NAME := status-badge
OS := $(shell uname)

build: clean
	rm -rf requirements.lock
	sed -i -e "s/version:.*/version: $(VERSION)/" Chart.yaml
	sed -i -e "s/imagetag: <VERSION>/imagetag: $(VERSION)/" values.yaml
	jx step helm build
	helm lint
	helm template .

install: clean build
	helm install . --name ${NAME}

upgrade: clean build
	helm upgrade ${NAME} .

delete:
	helm delete --purge ${NAME}

clean:
	rm -rf charts
	rm -rf ${NAME}*.tgz

release: clean build
	helm package .
	curl --fail -u $(CHARTMUSEUM_CREDS_USR):$(CHARTMUSEUM_CREDS_PSW) --data-binary "@$(NAME)-$(VERSION).tgz" $(CHART_REPO)/api/charts
	rm -rf ${NAME}*.tgz
