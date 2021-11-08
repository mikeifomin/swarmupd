# todo

 [ ] wait for coverage or rollback

# dev

docker stack deploy --compose-file ./docker-compose.yml

# example .gitlab-ci.yml


``` yaml
stages:  
  - build
  - deploy
deploy-prod-swarm:
  stage: deploy
  image:  mikeifomin/swarmupd:new
  script: 
    - >
      swarmupd update 
        --token $PROD_SWARMUPD_TOKEN
        --url $PROD_SWARMUPD_URL 
        --service-id midas_prod_api_public
        --new-tag $CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
docker-build:
  image: docker:latest
  stage: build
  services:
    - docker:dind
  before_script:
    - docker login -u "$CI_REGISTRY_USER" -p "$CI_REGISTRY_PASSWORD" $CI_REGISTRY
  # Default branch leaves tag empty (= latest tag)
  # All other branches are tagged with the escaped branch name (commit ref slug)
  script:
    - |
      if [[ "$CI_COMMIT_BRANCH" == "$CI_DEFAULT_BRANCH" ]]; then
        tag=""
        echo "Running on default branch '$CI_DEFAULT_BRANCH': tag = 'latest'"
      else
        tag=":$CI_COMMIT_BRANCH"
        echo "Running on branch '$CI_COMMIT_BRANCH': tag = $tag"
      fi
    - docker build --pull -t "$CI_COMMIT_SHA" --build-arg "COMMIT_SHA=$CI_COMMIT_SHORT_SHA" .
    - docker tag "$CI_COMMIT_SHA" "$CI_REGISTRY_IMAGE${tag}" 
    - docker tag "$CI_COMMIT_SHA" "$CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA"
    - docker push "$CI_REGISTRY_IMAGE${tag}"
    - docker push "$CI_REGISTRY_IMAGE:$CI_COMMIT_SHORT_SHA"
  rules:
    - if: $CI_COMMIT_BRANCH

```
