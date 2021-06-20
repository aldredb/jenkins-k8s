pipeline {
    agent any
    environment {
        registry = "jenkins-hello"
        ecrAddress = "878244765130.dkr.ecr.ap-southeast-1.amazonaws.com"
        awsCredential = 'aws'
        dockerImage = ''
        assumeRoleArn = 'arn:aws:iam::878244765130:role/EksWorkshopCodeBuildKubectlRole'
        clusterName = 'cluster-1'
    }
    stages {
         stage('Building image') {
            steps{
                script {
                    dockerImage = docker.build(registry + ":$BUILD_NUMBER")
                }
            }
        }
        stage('Push to ECR') {
            steps{
                script {
                    docker.withRegistry( "https://$ecrAddress/$registry", "ecr:ap-southeast-1:$awsCredential" ) {
                        dockerImage.push()
                    }
                }
            }
        }
        stage('Remove Unused docker image') {
            steps{
                sh "docker rmi $registry:$BUILD_NUMBER"
            }
        }
        stage('Modify K8S Manifest'){
            steps {
                sh '''
                    export IMAGE="$ecrAddress/$registry:$BUILD_NUMBER"
                    sed -ie "s~IMAGE~$IMAGE~g" hello.yaml
                   '''
            }
        }
        stage('Deploy to Cluster'){
            steps {
                withAWS(region:'ap-southeast-1',credentials:'aws') {
                    sh '''
                       CREDENTIALS=$(aws sts assume-role --role-arn $assumeRoleArn --role-session-name jenkins-kubectl --duration-seconds 900)
                       export AWS_ACCESS_KEY_ID="$(echo ${CREDENTIALS} | jq -r '.Credentials.AccessKeyId')"
                       export AWS_SECRET_ACCESS_KEY="$(echo ${CREDENTIALS} | jq -r '.Credentials.SecretAccessKey')"
                       export AWS_SESSION_TOKEN="$(echo ${CREDENTIALS} | jq -r '.Credentials.SessionToken')"
                       export AWS_EXPIRATION=$(echo ${CREDENTIALS} | jq -r '.Credentials.Expiration')
                       aws eks --region ap-southeast-1 update-kubeconfig --name $clusterName
                       kubectl apply -f hello.yaml -n jenkins
                       kubectl rollout status deploy/hello --timeout=600s -n jenkins
                       '''                    
                }
            }
        }
    }
}