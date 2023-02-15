run:
	go build . 
	./IMDb6dos 
clean:
	rm -r data/
	rm IMDb6dos
	mkdir data/
