pipeline {
    agent { docker { image 'docker:dind' } }
    stages {
        stage('build') {
            steps {
                script {
                    writeFile file: 'Dockerfile.build', text: """
                        FROM golang:1.22-alpine as builder
                        WORKDIR /app
                        COPY . .
                        RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./main .
                    """

                    sh "docker build -f Dockerfile.build -t build-image ."

                    archiveArtifacts artifacts: './main', fingerprint: true
                }
            }
        }
    }
}
