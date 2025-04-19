.PHONY: generate clean

# way to generator
GENERATOR = ./generator/generate_topics.go

# out file
OUTPUT = topics.go

# code get from yaml
generate:
	@echo "🔧 Generating $(OUTPUT) from topics.yaml..."
	go run $(GENERATOR)
	@echo "✅ Done."

# gen files cleaner
clean:
	@rm -f $(OUTPUT)
	@echo "🧹 Cleaned generated files."
