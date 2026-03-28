export const LogProfiler = (
  ProfilerId,
  Phase,
  ActualTime,
  Basetime,
  StartTime,
  CommitTime,
  Interactions,
) => {
  console.log({
    ProfilerId,
    Phase,
    ActualTime,
    Basetime, //time taken by react
    StartTime, //time at which render starts
    CommitTime,
    Interactions, // this is gotten from the rapping API
  });
};
export const RandomColour = () =>
  '#' + ((Math.random() * 0xffffff) << 0).toString(16);
