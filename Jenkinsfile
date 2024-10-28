pipeline {
    agent { docker { image 'golang:1.22-alpine' } }
    stages {
        stage('build') {
            steps {
                script {
                    sh "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./main ."
                    archiveArtifacts artifacts: './main', fingerprint: true
                }
            }
        }
    }
}
