package install

import (
	"github.com/kumahq/kuma/app/kumactl/pkg/install/data"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"reflect"
	"strings"
)

func renderHelmFiles(templates []data.File, args interface{}) ([]data.File, error) {

	files := []*loader.BufferedFile{}
	for _, f := range templates {
		files = append(files, &loader.BufferedFile{
			Name: f.FullPath[1:],
			Data: f.Data,
		})
	}

	loadedChart, err := loader.LoadFiles(files)
	if err != nil {
		return nil, err
	}

	loadedTemplates := loadedChart.Templates
	loadedChart.Templates = []*chart.File{}

	for _, t := range loadedTemplates {
		if !strings.HasPrefix(t.Name, "templates/pre-") {
			loadedChart.Templates = append(loadedChart.Templates, &chart.File{
				Name: t.Name,
				Data: t.Data,
			})
		}
	}

	overrideValues := map[string]interface{}{}

	v := reflect.ValueOf(args)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		value := v.FieldByName(name).Interface()

		splitName := strings.Split(name, "_")
		len := len(splitName)

		root := overrideValues

		for i := 0; i < len-1; i++ {
			n := splitName[i]

			if _, ok := root[n]; !ok {
				root[n] = map[string]interface{}{}
			}
			root = root[n].(map[string]interface{})
		}
		root[splitName[len-1]] = value
	}

	if err := chartutil.ProcessDependencies(loadedChart, overrideValues); err != nil {
		return nil, err
	}

	options := chartutil.ReleaseOptions{
		Name:      "kuma",
		Namespace: "kuma-system",
		Revision:  1,
		IsInstall: true,
		IsUpgrade: false,
	}
	valuesToRender, err := chartutil.ToRenderValues(loadedChart, overrideValues, options, nil)

	out, err := engine.Render(loadedChart, valuesToRender)
	if err != nil {
		return nil, errors.Errorf("Failed to render templates: %s", err)
	}
	result := []data.File{}

	for _, crd := range loadedChart.CRDObjects() {
		result = append(result, data.File{
			Data: crd.File.Data,
			Name: crd.Name,
		})
	}

	for n, d := range out {
		if strings.HasSuffix(n, "yaml") {
			result = append(result, data.File{
				Data: []byte(d),
				Name: n,
			})
		}
	}

	return result, nil
}
