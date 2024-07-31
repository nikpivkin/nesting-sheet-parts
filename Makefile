DATASETS = swim_2007-05-15.zip shirts_2007-05-15.zip trousers_2007-05-15.zip
DATASETS_URL = https://www.euro-online.org/websites/esicup/wp-content/uploads/sites/12/2017/04
DATASETS_FOLDER = datasets

download-datasets:
	@rm -rf $(DATASETS_FOLDER)
	@mkdir -p $(DATASETS_FOLDER)
	@for file in $(DATASETS); do wget -q $(DATASETS_URL)/$$file -P $(DATASETS_FOLDER); done
	@for file in $(DATASETS); do tar -xf $(DATASETS_FOLDER)/$$file -C $(DATASETS_FOLDER); done

run-swim:
	go run . --dataset $(DATASETS_FOLDER)/swim_2007-05-15/swim.xml --scale-output 0.04 --resolution 110

run-shirts:
	go run . --dataset $(DATASETS_FOLDER)/shirts_2007-05-15/shirts.xml --scale-output 5 --resolution 0.2 \
		--rotations 0

run-trousers:
	go run . --dataset $(DATASETS_FOLDER)/trousers_2007-05-15/trousers.xml --scale-output 1.5 --resolution 0.2