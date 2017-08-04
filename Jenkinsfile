node('golang') {
	def root = tool name: 'Go 1.8', type: 'go'
	def version = '0.1.5'
  def gitUrl = 'git@github.com:csduarte/mattermost-probe.git'
  def projectName = "mattermost-probe"
  def filename = "${applicationName}-${env.BUILD_NUMBER}"

  deleteDir()

  git([url: gitUrl, branch: env.BRANCH_NAME, credentialsId: 'uchat-mobile-key'])

	stage('prep') {
		withEnv(["GOROOT=${root}", "GOPATH=${WORKSPACE}", "PATH+GO=${root}/bin"]) {
			sh 'go version'
			sh 'cd $WORKSPACE/src/github.com/csduarte/mattermost-probe && glide install'
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