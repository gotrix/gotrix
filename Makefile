includeDir := "include"
componentsDir := "components"

generate:
	go generate ./...

# build components in ./components
components:
	@for comp in $$(find ${componentsDir} -print | grep ".go" | sed "s/${componentsDir}\///g"); do \
  		echo "building component $${comp}"; \
		go build -buildmode=plugin -o ./${componentsDir}/$$(echo "$${comp}" | sed "s/\.go/\.so/g") ${componentsDir}/$${comp}; \
	done

.PHONY: components generate