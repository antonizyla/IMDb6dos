run:
	go build . 
	./IMDb6dos
seed:
	go build . 
	./IMDb6dos --seed
clean:
	rm -r data/
	rm IMDb6dos
	mkdir data/
