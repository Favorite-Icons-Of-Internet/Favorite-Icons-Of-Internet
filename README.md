# Favorite Icons Of Internet

Art project that aims to depict the vastness and colorfullness of the internet.

You can see the result of all the crawling and image-crunching at [FavoriteIconsOfInternet.com](http://www.favoriteiconsofinternet.com/)

Our current goal is to bring the project to the state where we can keep the history of daily favicon changes for at least a million web sites.

![Favorite Icons](favicons-project-illustration.png)

## Processing Steps

### Step 1. [Load domains](https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/2) (on central box)

Updates a list of domains in the database, currently takes a list of Alexa Rankings.

### Step 2. [Get a list of domains to crawl](https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/3) (on central box)

Gets a list of domains to crawl (currently only active Alexa domains) and uploads them to a queue in chunks for crawlers to pick up

### Step 3. [Fetch icons](https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/1) (on crawler workers)

Listens for messages in a queue and crawls the sites in the message finding favorite icons and comparing them to existing version to see if the have changed.

### Step 4. [Convert icons to PNH](https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/4) (on crawler workers)

After all icons are fetched, convert them to PNG, calculate average color and upload to results storage together with manifest describing which icons are new, which has changed and etc.

### Step 5. [Calculate tiles to be updated] (https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/4) (on central box)

Gather all the results and update the database. Calculate a list of tiles that need to be updated (currently all tiles with predefined width/height ordered by Alexa ranking) and put each tile as a job into a queue.

Generate HTML and necessary JSON metadata/

### Step 6. [Generate tiles](https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/6) (on tile workers)

Grab images required for the tile (or sync them all) and generate a tile. Optimize the image using smu.sh and deploy to a CDN.

### Step 7. [Move HTML and metadata to production](https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/6) (on central box)

Once all tiles are done, move HTML and metadata chunks over to production!

## Step 8. [Send emails, daily reports and etc](https://github.com/Favorite-Icons-Of-Internet/Favorite-Icons-Of-Internet/issues/8) (on central box)

Notify users (if any), send daily newsletter and etc.
