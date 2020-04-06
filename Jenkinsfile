pipeline{
    agent any
    environment {
        registry = "2017330017/go-user"
        GOCACHE = "/tmp"
    }
    stages {
        stage('Build') {
            agent {
                docker {
                    image 'golang'
                }
            }
            steps {
                // create project directory
                sh 'cd ${GOPATH}/src'
                sh 'mkdir -p ${GOPATH}/src/guser'
                // copy all files in our Jenkins workspace to our project directory
                sh 'cp -r ${WORKSPACE}/* ${GOPATH}/src/guser'
                // build the app
                sh 'go build'
            }
        }
        stage('Publish') {
            environment {
                registryCredential = 'dockerhub'
            }
            steps {
                script {
                    def appimage = docker.build registry + ":$BUILD_NUMBER"
                    docker.withRegistry('', registryCredential) {
                        appimage.push()
                        appimage.push('latest')
                    }
                }
            }
        }
        stage('Deploy') {
            steps {
                script {
                    def image_id = registry + ":$BUILD_NUMBER"
                    sh "ansible-playbook playbook.yml --extra-vars \"image_id=${image_id}\""
                }
            }
        }
    }
}