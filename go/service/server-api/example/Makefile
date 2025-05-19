# directories
ROOT_DIR   := $(CURDIR)
OUTPUT_DIR := $(ROOT_DIR)/build

# cmd
MKDIR = mkdir -p
ECHO  = $(if $(filter Linux,$(shell uname)),echo -e,echo)
CD    = cd
CP    = cp
RM    = rm

all:
	$(MKDIR) $(OUTPUT_DIR)

	@$(ECHO) "Building server example..."
	CGO_ENABLED=0 go build -o $(OUTPUT_DIR)/server_example $(ROOT_DIR)/main.go
	@$(ECHO) "Built successfully: $(OUTPUT_DIR)/server_example"

clean:
	@$(ECHO) "Cleaning..."
	$(RM) -rf $(ROOT_DIR)/build
	@$(ECHO) "Cleaned successfully"