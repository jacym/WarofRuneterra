pipeline {
  agent {
    docker {
      image 'golang:alpine'
    }

  }
  stages {
    stage('Pull Dependencies') {
      agent any
      steps {
        dir(path: 'server') {
          sh 'go get'
        }

      }
    }

    stage('Build It') {
      parallel {
        stage('Lint') {
          agent any
          steps {
            dir(path: 'server') {
              sh 'go vet'
            }

          }
        }

        stage('Build') {
          agent any
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