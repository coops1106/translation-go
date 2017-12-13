.PHONY: generateTranslationFiles


# go get -u github.com/nicksnyder/go-i18n/goi18n
# default sourceLanguage is en-us
# sourceLanguage acts as the seed file and therefore default values
generateTranslationFiles:
	rm -rf ./data/output
	mkdir -p ./data/output/gen
	mkdir ./data/output/final
	mkdir ./data/output/constants

	goi18n -sourceLanguage en-gb -outdir ./data/output/gen ./data/input/*.json
	goi18n -outdir ./data/output/final -format toml ./data/output/gen/*.all.json ./data/output/gen/*.untranslated.json
	goi18n constants -outdir ./data/output/constants -package constants ./data/output/final/en-gb.all.toml