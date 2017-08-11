node('golang') {
	def root			= tool name: 'Go 1.8', type: 'go'
	def gitUrl		= 'git@github.com:csduarte/mattermost-probe.git'
	def project		= "mattermost-probe"
	def bucket		= "uchat-releases"
	def date			= new Date().format( 'yyMMdd' )
	def filename	= "${project}-${date}-${env.BUILD_NUMBER}"

  deleteDir()

	checkout([$class: 'GitSCM',
		branches: [[name: '*/master']], 
		doGenerateSubmoduleConfigurations: false, 
		extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'src/github.com/csduarte/mattermost-probe']], 
		submoduleCfg: [], 
		userRemoteConfigs: [[credentialsId: 'uchat-mobile-key',url: gitUrl]]])

	stage('prep') {
		withEnv(["GOROOT=${root}", "GOPATH=${WORKSPACE}", "PATH+GO=${root}/bin"]) {
			sh 'go version'
			sh "cd $WORKSPACE/src/github.com/csduarte/mattermost-probe && $JENKINS_HOME/go/bin/glide install"
			sh "cd $WORKSPACE/src/github.com/csduarte/mattermost-probe && make .prebuild"
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
			sh 'cd $WORKSPACE && go build github.com/csduarte/mattermost-probe'
			sh "cp mattermost-probe ${filename}"
		}
	}

	stage('publish') {
    archiveArtifacts artifacts: filename, fingerprint: true
      
    withAWS(credentials: 'aws-uchat-releases', region: 'us-west-2') {
      s3Upload  bucket: bucket,
                file: 	filename, 
                path: 	"${project}/${env.BRANCH_NAME}/${filename}"
		} 
	}
}
