pipeline {
  agent {
    docker {
      image 'golang:alpine'
    }

  }
  stages {
    stage('Pull Dependencies') {
      steps {
        dir(path: 'server') {
          sh 'go get'
        }

      }
    }

    stage('Lint') {
      parallel {
        stage('Lint') {
          agent {
            docker {
              image 'golang:alpine'
            }

          }
          steps {
            dir(path: 'server') {
              sh 'go vet'
            }

          }
        }

        stage('Build') {
          agent {
            docker {
              image 'golang:alpine'
            }

          }
          steps {
            dir(path: 'server') {
              sh 'go build'
            }

          }
        }

      }
    }

  }
}