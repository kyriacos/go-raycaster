# go-raycaster

Simple raycaster in Go using SDL.

There's plenty of comments in the code. Check it out if you are interested and just email me if you have any questions.

_The textures in the images directory are from Wolfenstein 3D and all copyrights belong to ID Software._

## Notes:

#### Rendering - FPS

Through trial and error mostly i got the Raycaster running from an average of 10fps to an average of 48fps.

As a first pass i went through all the instances of images that i had and made sure everything was in the same format i.e. NRGBA. Then i went through and i replaced all the calls to `.Set` with `.SetNRGBA`. Those two changes alone bumped the fps up to 20!

It seems that there are a lot of interface checks or type assertions i suppose on the `.Set` method and using the `.SetNRGBA` method is A LOT faster. To make sure that was an issue i simply changed a single call from `.SetNRGBA` to just `.Set` and the fps dropped by 7fps.
I also tried using the `draw.Image` interface instead but that didn't do anything.

Eventually i switched to having a fixed sized `uint32` array (`[WindowWidth*WindowHeight]uint32`) which i converted to a `[]byte` slice every time i updated the texture in SDL. Funny enough that gave me another 10fps and got the total average up to 30fps!

As a last pass i just implemented everything using a `[]byte` slice and added some convenience methods. In fact it actually resembles the `image.NRGBA` struct. That brought up the average FPS to 40 i believe.

And the last optimization i made was removing the call to `.Clear()`. I can't tell a difference. To be honest although all the tutorials say you should clear out the buffer every time, in my scenario the entire buffer is re-rendered every time so i skipped it :). Now it went up to an average of 50fps. I'll stop there. Let me know if there are more things i could do or if you think there is something wrong with what i am doing.

I am very curious as to why the image.NRGBA was so much slower. I actually went back and just used an image.NRGBA again just to make sure i wasn't doing anything wrong and the FPS dropped again.

## TODO:

- [x] Load levels from external file
- [ ] Remove extra global VARS
- [ ] Determine texture size dynamically or set as a constant
- [ ] Add textures for floor and ceiling
- [ ] Clean up FPS calculation code
- [ ] Don't use the global gameMap variable in Ray.cast
- [x] Optimize the code. Figure out a way to make it run faster. It's very slow...
