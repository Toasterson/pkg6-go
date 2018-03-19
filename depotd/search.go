package depotd

import (
	"github.com/labstack/echo"
)

/*
The main difference between search 0 and 1 seems to be that 0 reports all the details in the request while 1 leaves some details to be read by the client from disk.
*/

func (d *DepotServer)handleSearchV0(c echo.Context) error {
	return nil
}

/*
/{publisher}/search/0/{query}
query is p<:::vim> where gnu is the search string
pkg.summary pkg:/library/augeas-vim@1.4.0,5.11-2017.0.0.1:20170306T182239Z set Augeas-Tools is a configuration-file editing C library and shell toolkit - ViM editor binding
pkg.description pkg:/editor/vim@8.0.104,5.11-2017.0.0.2:20170306T113401Z set Vim is a clone of the Unix editor 'vi'.  It is a modal text editor with support for syntax highlighting, context-sensitive indentation, and extension scripting in numerous languages.
pkg.description pkg:/editor/vim@8.0.104,5.11-2017.0.0.2:20170609T180153Z set Vim is a clone of the Unix editor 'vi'.  It is a modal text editor with support for syntax highlighting, context-sensitive indentation, and extension scripting in numerous languages.
basename pkg:/editor/vim/vim-core@8.0.104,5.11-2017.0.0.2:20170306T113419Z file usr/bin/vim
basename pkg:/editor/vim/vim-core@8.0.104,5.11-2017.0.0.2:20170609T180211Z file usr/bin/vim
com.oracle.info.name pkg:/editor/gvim@8.0.104,5.11-2017.0.0.2:20170609T180205Z set vim
com.oracle.info.name pkg:/editor/vim/vim-core@8.0.104,5.11-2017.0.0.2:20170306T113419Z set vim
com.oracle.info.name pkg:/editor/vim/vim-core@8.0.104,5.11-2017.0.0.2:20170609T180211Z set vim
com.oracle.info.name pkg:/editor/vim@8.0.104,5.11-2017.0.0.2:20170306T113401Z set vim
com.oracle.info.name pkg:/editor/gvim@8.0.104,5.11-2017.0.0.2:20170306T113411Z set vim
com.oracle.info.name pkg:/editor/vim@8.0.104,5.11-2017.0.0.2:20170609T180153Z set vim
pkg.description pkg:/editor/vim/vim-core@8.0.104,5.11-2017.0.0.2:20170306T113419Z set This package delivers the core executables and man pages for vim (pkg:/editor/vim), and will generally only be installed independently on a minimized system.
pkg.description pkg:/editor/vim/vim-core@8.0.104,5.11-2017.0.0.2:20170609T180211Z set This package delivers the core executables and man pages for vim (pkg:/editor/vim), and will generally only be installed independently on a minimized system.
pkg.fmri pkg:/editor/vim@8.0.104,5.11-2017.0.0.2:20170306T113401Z set openindiana.org/editor/vim
pkg.fmri pkg:/editor/vim@8.0.104,5.11-2017.0.0.2:20170609T180153Z set openindiana.org/editor/vim
*/

func (d *DepotServer)handleSearchV1(c echo.Context) error {
	return nil
}

/*
/{publisher}/search/1/{casesensitive}_{returntype}_{maxreturn}_{startreturn}_p<{query}>
returntype can be 1 or 2 but does nothing
query is p<:::gnu> where gnu is the search string
query format pkg_name:action_name:index:token
pkg_name being the package name
action_name being the type of action to search e.g. dir,file etc.
index being the attribute of the action to search
token being the search string wildcads like * and ? are supported.
casesensitive is for switching casesensitivity on and of
maxreturn and startreturn are for pageing
/openindiana.org/search/1/False_1_None_None_p%3C%3A%3A%3Agnupg%3E
Return from search v1
0 1 pkg:/crypto/gnupg@2.0.30,5.11-2017.0.0.0:20170306T211520Z
0 1 pkg:/crypto/gnupg@2.0.30,5.11-2017.0.0.0:20170810T075528Z
0 1 pkg:/crypto/gnupg@2.0.30,5.11-2017.0.0.0:20170810T105136Z
0 1 pkg:/library/security/gpgme@1.1.8,5.11-2017.0.0.1:20170307T004012Z
0 1 pkg:/library/security/gpgme@1.1.8,5.11-2017.0.0.1:20170810T075231Z
0 1 pkg:/library/security/gpgme@1.1.8,5.11-2017.0.0.1:20170810T095659Z
0 1 pkg:/library/security/libassuan@2.1.3,5.11-2017.0.0.0:20170306T225546Z
0 1 pkg:/library/security/libassuan@2.1.3,5.11-2017.0.0.1:20170606T194953Z
0 1 pkg:/library/security/libassuan@2.1.3,5.11-2017.0.0.1:20170810T074932Z
0 1 pkg:/library/security/libgpg-error@1.19,5.11-2017.0.0.0:20170306T124025Z
0 1 pkg:/library/security/libgpg-error@1.27,5.11-2017.0.0.0:20170810T071845Z
0 1 pkg:/library/security/libksba@1.3.4,5.11-2017.0.0.0:20170306T230029Z
0 1 pkg:/library/security/libksba@1.3.5,5.11-2017.0.0.0:20170606T133035Z
0 1 pkg:/library/security/libksba@1.3.5,5.11-2017.0.0.0:20170810T074723Z
0 1 pkg:/security/pinentry@0.9.7,5.11-2017.0.0.0:20170307T004027Z
0 1 pkg:/security/pinentry@0.9.7,5.11-2017.0.0.0:20170810T075406Z
0 1 pkg:/system/library/security/libgcrypt@1.5.6,5.11-2017.0.0.0:20170306T133652Z
0 1 pkg:/system/library/security/libgcrypt@1.5.6,5.11-2017.0.0.1:20170714T085321Z
0 1 pkg:/system/library/security/libgcrypt@1.5.6,5.11-2017.0.0.1:20170810T082341Z
0 1 pkg:/system/library/security/libgcrypt@1.8.0,5.11-2017.0.0.0:20170810T084351Z
0 1 pkg:/system/library/security/libgcrypt@1.8.1,5.11-2017.0.0.0:20170913T183858Z
*/
