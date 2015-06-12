	REXIntro = regexp.MustCompile("(/\\*\\*Intro[\\s]+(\\s|.)*?[\\s]+\\*/)|(//Intro[\\s]?([\\s]|.)*?//[Ee][Nn][Dd])")

	// parse public document content
	/** *[z-index-][title]
	 *  document text or markdown text
	 */
	// /[z-index-][title]
	//  document text or markdown text
