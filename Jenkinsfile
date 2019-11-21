pipeline {
  agent {
    docker {
      image 'golang:alpine'
    }

  }
  stages {
    stage('Pull Dependencies') {
      steps {
        dir(path: 'server')
        sh '''# grab dependencies
go get'''
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
            dir(path: 'server')
            sh '''# lint
go vet'''
          }
        }

        stage('Build') {
          steps {
            dir(path: 'server')
            sh '''# build
go build'''
          }
        }

      }
    }

  }
}