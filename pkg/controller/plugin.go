package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/araminian/argo-appset-pr-label-filter/pkg/models"
	"github.com/araminian/argo-appset-pr-label-filter/pkg/utils"
)

func FilterPRs(w http.ResponseWriter, r *http.Request) {

	token, err := utils.ReadToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		fmt.Printf("Error reading token: %s\n", err.Error())
		return
	}

	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader != fmt.Sprintf("Bearer %s", token) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		fmt.Println("Unauthorized")
		return
	}

	// read the body
	PluginInputs := &models.PluginInputs{}
	utils.ParseBody(r, PluginInputs)

	// print the labels
	labels := PluginInputs.Input.Parameters.GetLabels()
	fmt.Println("labels: ", labels)
	// check if the preview is deployable
	deployable := PluginInputs.Input.Parameters.Deployable()

	if deployable {
		w.WriteHeader(http.StatusOK)
		fmt.Println("ApplicationSet: ", PluginInputs.ApplicationSetName)
		fmt.Println("PR Number: ", PluginInputs.Input.Parameters.Number)
		fmt.Println("Exclude label not found")
		fmt.Println("Using Path: ", PluginInputs.Input.Parameters.Path)

		output := &models.PluginOutputs{}

		newOutputParams := []models.OutputsParameters{
			{GeneratedPath: PluginInputs.Input.Parameters.Path},
		}

		output.SetOutputParameters(newOutputParams)

		fmt.Println("output: ", output)

		res, _ := json.Marshal(output)

		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	}

	fmt.Println("ApplicationSet: ", PluginInputs.ApplicationSetName)
	fmt.Println("PR Number: ", PluginInputs.Input.Parameters.Number)
	fmt.Println("Exclude label found")
	fmt.Println("Using BlackHolePath: ", PluginInputs.Input.Parameters.BlackHole)

	output := &models.PluginOutputs{}

	newOutputParams := []models.OutputsParameters{
		{GeneratedPath: PluginInputs.Input.Parameters.BlackHole},
	}

	output.SetOutputParameters(newOutputParams)

	fmt.Println("output: ", output)

	res, _ := json.Marshal(output)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
