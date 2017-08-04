node('golang') {
	def root = tool name: 'Go 1.8', type: 'go'
  def gitUrl = 'git@github.com:csduarte/mattermost-probe.git'
  def projectName = "mattermost-probe"
	def date = new Date().format( 'yyMMdd' )
  def filename = "${projectName}-${date}-${env.BUILD_NUMBER}"

  deleteDir()

  git([url: gitUrl, branch: env.BRANCH_NAME, credentialsId: 'uchat-mobile-key'])
	checkout([$class: 'GitSCM',
		branches: [[name: '*/master']], 
		doGenerateSubmoduleConfigurations: false, 
		extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'src/github.com/csduarte/mattermost-probe']], 
		submoduleCfg: [], 
		userRemoteConfigs: [[credentialsId: 'uchat-mobile-key',url: gitUrl]]])

	stage('prep') {
		withEnv(["GOROOT=${root}", "GOPATH=${WORKSPACE}", "PATH+GO=${root}/bin"]) {
			sh 'go version'
			sh 'cd $WORKSPACE/src/github.com/csduarte/mattermost-probe && $PATH+GO/glide install'
			sh 'if [[ ! -d $WORKSPACE/bin ]]; then mkdir $WORKSPACE/bin; fi; if [[ ! -d $WORKSPACE/pkg ]]; then mkdir $WORKSPACE/pkg; fi'
		}
	}

	stage('test') {
		withEnv(["GOROOT=${root}", "GOPATH=${WORKSPACE}", "PATH+GO=${root}/bin"]) {
			sh 'cd $WORKSPACE && go test github.com/csduarte/mattermost-probe/mattermost  '
		}
	}

	stage('build') {
		withEnv(["GOROOT=${root}", "GOPATH=${WORKSPACE}", "PATH+GO=${root}/bin"]) {
			sh 'cd $WORKSPACE && go build -o ${filename} github.com/csduarte/mattermost-probe'
		}
	}

	stage('publish') {
    archiveArtifacts artifacts: filename, fingerprint: true
      
    withAWS(credentials: 'aws-uchat-releases', region: 'us-west-2') {
      s3Upload  bucket: 'uchat-releases',
                file: filename, 
                path: "${projectName}/${env.BRANCH_NAME}/${filename}"
		} 
	}
}
