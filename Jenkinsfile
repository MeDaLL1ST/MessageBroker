pipeline {
    agent { docker { image 'golang:1.22-alpine' } }
    stages {
        stage('build') {
            steps {
                script {
                    sh "go build -o ./main ."
                    archiveArtifacts artifacts: './main', fingerprint: true
                }
            }
        }
    }
}
