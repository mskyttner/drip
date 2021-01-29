#* @get /counter
function(req){
  count <- 0
  if (!is.null(req$session$counter)){
    count <- as.numeric(req$session$counter)
  }
  req$session$counter <- count + 1
  return(paste0("This is visit #", count))
}

#* @get /drip/version
function() {
	system("drip version", intern = TRUE)
}

#* @get /drip/routes
function() {
	system("drip routes", intern = TRUE)
}

#* @get /drip/mountdir
function() {
	Sys.getenv("PLUMBER_DIR")
}

#* @get /drip/mountdir/files
function() {
	dir(Sys.getenv("PLUMBER_DIR"))
}

#* @get /plumber/apis/available
function() {
	plumber::available_apis()
}

#* @get /plumber/version
function() {
	as.character(packageDescription("plumber"))
}

#* @get /sessioninfo
function() {
	as.character(sessionInfo())
}

#* @get /readme
function(res) {
   res$status <- 302
   res$setHeader("Location", "/__docs__/")
}

