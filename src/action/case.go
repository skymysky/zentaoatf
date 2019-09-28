package action

import (
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
)

func CommitCases(files []string) {
	for _, file := range files {
		pass, id, _, title := zentaoUtils.GetCaseInfo(file)

		if pass {
			cpStepDescArr, cpStepDescArr, cpExpectArr := zentaoUtils.ReadScriptCheckpoints(file)

			zentaoService.CommitCase(id, title)
		}
	}
}
