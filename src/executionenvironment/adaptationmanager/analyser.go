package adaptationmanager

import (
	"framework/configuration/configuration"
	"shared/shared"
	"shared/parameters"
	"fmt"
)

type Analyser struct{}

func (Analyser) Exec(conf configuration.Configuration, chanMACorrective chan shared.MonitoredCorrectiveData, chanMAEvolutive chan shared.MonitoredEvolutiveData, chanMAProactive chan shared.MonitoredProactiveData, chanAP chan shared.AnalysisResult) {

	// prepapre channels
	chanCorrective := make(chan shared.AnalysisResult)
	chanEvolutive := make(chan shared.AnalysisResult)
	chanProactive := make(chan shared.AnalysisResult)

	if parameters.IS_CORRECTIVE {
		go correctiveAnalysis(chanMACorrective, chanCorrective)
	}
	if parameters.IS_EVOLUTIVE {
		go evolutiveAnalysis(chanMAEvolutive, chanEvolutive)
	}
	if parameters.IS_PROACTIVE {
		go proactiveAnalysis(chanMAProactive, chanProactive)
	}

	for {
		select {
		case analysisResult := <-chanCorrective:
			chanAP <- analysisResult
		case analysisResult := <-chanEvolutive:
			fmt.Println("Analyser:: Evolutive")
			chanAP <- analysisResult
		case analysisResult := <-chanProactive:
			chanAP <- analysisResult
		}
	}
}

func correctiveAnalysis(chanMa chan shared.MonitoredCorrectiveData, chanCorrective chan shared.AnalysisResult) {

	for {
		monitoredData := <-chanMa
		r := invokePROM(monitoredData)
		if r {
			chanCorrective <- shared.AnalysisResult{Analysis: parameters.NO_CHANGE} // TODO
		}
	}
}

func proactiveAnalysis(chanMa chan shared.MonitoredProactiveData, chanProactive chan shared.AnalysisResult) { // TODO
	for {
		monitoredData := <-chanMa
		r := invokePRISM(monitoredData)
		if r {
			chanProactive <- shared.AnalysisResult{Analysis: parameters.NO_CHANGE}
		}
	}
}

func evolutiveAnalysis(chanMa chan shared.MonitoredEvolutiveData, chanEvolutive chan shared.AnalysisResult) {
	analysisResult := shared.AnalysisResult{}

	for {
		listOfNewPlugins := <-chanMa // receive new plugins from Monitor
		analysisResult.Result = listOfNewPlugins
		analysisResult.Analysis = parameters.EVOLUTIVE_CHANGE
		chanEvolutive <- analysisResult
	}
}

func invokePROM(data shared.MonitoredCorrectiveData) bool {
	// TODO
	return false
}

func invokePRISM(data shared.MonitoredProactiveData) bool {
	// TODO
	return false
}
