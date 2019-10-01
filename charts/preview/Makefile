preview: 
	sed -i -e "s/version:.*/version: $(VERSION)/" Chart.yaml
	sed -i -e "s/version:.*/version: $(VERSION)/" ../*/Chart.yaml
	sed -i -e "s/imagetag:.*/imagetag: $(VERSION)/" values.yaml
	sed -i -e "s/previewname: .*/previewname: status-badge-pr-$(PULL_NUMBER)/" values.yaml
	sed -i -e "s/fabric8\.io\/exposeUrl: .*/fabric8\.io\/exposeUrl: https:\/\/status-badge-pr-$(PULL_NUMBER).jenkins-x.live/" values.yaml
	echo "  version: $(VERSION)" >> requirements.yaml
	jx step helm build
