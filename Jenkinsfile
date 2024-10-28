pipeline {
    agent { docker { image 'golang:1.22-alpine' } }
    environment {
        HOME = "${env.WORKSPACE}"
    }
    stages {
        stage('build') {
            steps {
                script {
                    sh "go mod tidy"
                    sh "CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./main ."
                    archiveArtifacts artifacts: './main', fingerprint: true
                }
            }
        }
    }
}
